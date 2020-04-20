package com.lbg.kafka.opa;

import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import com.google.gson.Gson;
import kafka.security.auth.Acl;
import kafka.security.auth.Authorizer;
import kafka.security.auth.Operation;
import kafka.security.auth.Resource;
import lombok.Cleanup;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.common.security.auth.KafkaPrincipal;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static kafka.network.RequestChannel.Session;

@Slf4j
public class OpaAuthorizer implements Authorizer {

  private final static String OPA_AUTHORIZER_URL_CONFIG = "opa.authorizer.url";
  private final static String OPA_AUTHORIZER_DENY_ON_ERROR_CONFIG = "opa.authorizer.allow.on.error";
  private final static String OPA_AUTHORIZER_CACHE_INITIAL_CAPACITY_CONFIG = "opa.authorizer.cache.initial.capacity";
  private final static String OPA_AUTHORIZER_CACHE_MAXIMUM_SIZE_CONFIG = "opa.authorizer.cache.maximum.size";
  private final static String OPA_AUTHORIZER_CACHE_EXPIRE_AFTER_MS_CONFIG = "opa.authorizer.cache.expire.after.ms";
  private final static String OPA_AUTHORIZER_TOKEN = "opa.authorizer.token";

  private String opaUrl;
  private boolean allowOnError;
  private int initialCapacity;
  private int maximumSize;
  private long expireAfterMs;
  private String opaToken;

  private final Gson gson = new Gson();

  private final Map<String, Object> configs = new HashMap<>();

  private LoadingCache<Msg.Input, Boolean> cache;

  public OpaAuthorizer() {
    configure(new HashMap<>());
  }

  private LoadingCache<Msg.Input, Boolean> buildCache() {
    return CacheBuilder.newBuilder()
      .initialCapacity(initialCapacity)
      .maximumSize(maximumSize)
      .expireAfterWrite(expireAfterMs, TimeUnit.MILLISECONDS)
      .build(
        new CacheLoader<Msg.Input, Boolean>() {
          @Override
          public Boolean load(Msg.Input input) {
            return allow(input);
          }
        }
      );
  }

  private boolean allow(Msg.Input input) {
    try {
      HttpURLConnection conn = (HttpURLConnection) new URL(opaUrl).openConnection();

      conn.setDoOutput(true);
      conn.setRequestMethod("POST");
      conn.setRequestProperty("Content-Type", "application/json");
      if (!opaToken.isEmpty()) {
        conn.setRequestProperty("Authorization", "Bearer " + opaToken);
      }

      String data = gson.toJson(new Msg(input));
      OutputStream os = conn.getOutputStream();
      os.write(data.getBytes());
      os.flush();

      if (log.isTraceEnabled()) {
        log.trace("Response code: {}, Request data: {}", conn.getResponseCode(), data);
      }

      @Cleanup BufferedReader br = new BufferedReader(new InputStreamReader((conn.getInputStream())));
      return (boolean) gson.fromJson(br.readLine(), Map.class).get("result");
    } catch (IOException e) {
      return allowOnError;
    }
  }

  public String getValueOrDefault(String property, String defaultValue) {
    return Optional.ofNullable((String) configs.get(property)).orElse(defaultValue);
  }

  public boolean authorize(Session session, Operation operation, Resource resource) {
    try {
      return cache.get(new Msg.Input(operation, resource, session));
    } catch (ExecutionException e) {
      return allowOnError;
    }
  }

  public void configure(Map<String, ?> configs) {
    this.configs.clear();
    this.configs.putAll(configs);

    if (log.isTraceEnabled()) {
      log.trace("CONFIGS: {}", this.configs);
    }
    opaUrl = (String) getValueOrDefault(OPA_AUTHORIZER_URL_CONFIG, "http://localhost:8181");
    allowOnError = Boolean.valueOf((String) getValueOrDefault(OPA_AUTHORIZER_DENY_ON_ERROR_CONFIG, "false"));
    initialCapacity = Integer.parseInt((String) getValueOrDefault(OPA_AUTHORIZER_CACHE_INITIAL_CAPACITY_CONFIG, "100"));
    maximumSize = Integer.parseInt((String) getValueOrDefault(OPA_AUTHORIZER_CACHE_MAXIMUM_SIZE_CONFIG, "100"));
    expireAfterMs = Long.parseLong((String) getValueOrDefault(OPA_AUTHORIZER_CACHE_EXPIRE_AFTER_MS_CONFIG, "600000"));
    opaToken = (String) getValueOrDefault(OPA_AUTHORIZER_TOKEN, "");

    cache = buildCache();
  }

  public void addAcls(scala.collection.immutable.Set<Acl> acls, Resource resource) {
  }

  public boolean removeAcls(scala.collection.immutable.Set<Acl> acls, Resource resource) {
    return false;
  }

  public boolean removeAcls(Resource resource) {
    return false;
  }

  public scala.collection.immutable.Set<Acl> getAcls(Resource resource) {
    return null;
  }

  public scala.collection.immutable.Map<Resource, scala.collection.immutable.Set<Acl>> getAcls(KafkaPrincipal principal) {
    return null;
  }

  public scala.collection.immutable.Map<Resource, scala.collection.immutable.Set<Acl>> getAcls() {
    return null;
  }

  public void close() {
  }

  @Data
  static class Msg {
    private final Input input;

    @Data
    static class Input {
      private final Operation operation;
      private final Resource resource;
      private final Session session;
    }
  }
}

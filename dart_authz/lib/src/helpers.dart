import 'dart:io';

// Python-like os.getenv() implementation
String getenv(String envVar, [String defaultValue=""]) {
  Map<String, String> envVars = Platform.environment;
  return envVars[envVar] ?? defaultValue;
}

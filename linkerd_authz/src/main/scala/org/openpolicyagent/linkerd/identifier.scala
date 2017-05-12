package org.openpolicyagent.linkerd

import com.fasterxml.jackson.annotation.JsonIgnore
import com.twitter.finagle._
import com.twitter.finagle.http.{HeaderMap, Request, Response}
import com.twitter.io.Buf
import com.twitter.io.Buf.Utf8
import com.twitter.util.Future
import io.buoyant.linkerd.IdentifierInitializer
import io.buoyant.linkerd.protocol.HttpIdentifierConfig
import io.buoyant.linkerd.protocol.http.HeaderTokenIdentifierConfig
import io.buoyant.router.RoutingFactory
import io.buoyant.router.RoutingFactory.{Identifier, RequestIdentification, UnidentifiedRequest}
import io.buoyant.router.http.HeaderIdentifier
import org.json4s.DefaultFormats
import org.json4s.JsonAST.JValue
import org.json4s.jackson.{Serialization, parseJson}

case class AuthzIdentifier(
  prefix: Path,
  opaIp: String,
  opaPort: Int,
  opaPath: String,
  header: Option[String],
  baseDtab: () => Dtab = () => Dtab.base
) extends RoutingFactory.Identifier[Request] {

  import Helpers._

  def apply(req: Request): Future[RequestIdentification[Request]] = {

    // TODO(tsandall): support other built-in identifiers through config
    val identifier = HeaderIdentifier(prefix, header.getOrElse(HeaderTokenIdentifierConfig.defaultHeader), headerPath = false, baseDtab)

    // Execute normal identifier.
    identifier(req).flatMap { id =>

      // Prepare OPA query
      val addr = s"$opaIp:$opaPort"
      val client = Http.newService(addr)
      val request = http.Request(http.Method.Post, opaPath)

      request.host = addr // Finagle was not setting host automatically.
      request.setContentTypeJson()

      val pair = mapHeadersToSourceIPHostPair(req.headerMap) match {
        case Some((s, h)) => (Some(s), Some(h))
        case None => (None, None)
      }

      request.content = mapRequestToBuf(OPARequest(Input(req.method.toString(), req.path, req.headerMap, pair._1.map(Identity), pair._2)))

      // Execute OPA query
      client(request).map { response =>
        mapResponseToErrorMsg(response) match {
          case Some(msg) => new UnidentifiedRequest[Request](msg)
          case None => id
        }
      }
    }
  }
}

class AuthzIdentifierConfig(ip: String, port: Int, path: String, header: Option[String] = None) extends HttpIdentifierConfig{

  @JsonIgnore
  override def newIdentifier(prefix: Path, baseDtab: () => Dtab): Identifier[Request] = {
    new AuthzIdentifier(prefix, ip, port, path, header, baseDtab)
  }
}


class AuthzIdentifierInitializer extends IdentifierInitializer {

  override def configId: String = "org.openpolicyagent.linkerd.authzIdentifier"

  override def configClass: Class[_] = return classOf[AuthzIdentifierConfig]
}

case class OPARequest(input: Input)

case class Input(
                  method: String,
                  path: String,
                  headers: HeaderMap,
                  identity: Option[Identity],
                  host: Option[String]
                )

case class Identity(source_ip: String)

case class OPAResponse(result: Option[Document])

case class Document(errors: Seq[String] = Seq.empty[String])

case class OPAError(
                     code: String,
                     message: String,
                     errors: Option[Seq[Map[String,JValue]]]
                   ) {

  override def toString: String = s"$message (code: $code)"

}


private[linkerd] object Helpers {

  implicit val f = DefaultFormats

  def mapHeadersToSourceIPHostPair(headers: HeaderMap): Option[(String, String)] = {
    headers.get("Forwarded").flatMap({ s =>
      val parts = s.split(';')
      val items = parts.flatMap({ t =>
        val pair = t.split('=').toList match {
          case k :: v :: Nil => Some((k, v))
          case _ => None
        }
        pair match {
          case Some((k, v)) => List((k,v))
          case None => List.empty[(String, String)]
        }
      }).toMap
      (items.get("for"), items.get("host")) match {
        case (Some(sourceIP), Some(host)) => Some((sourceIP, host))
        case _ => None
      }
    })
  }

  def mapRequestToBuf(request: OPARequest): Buf = {
    val s = Serialization.write(request)
    Utf8(s)
  }

  def mapResponseToErrorMsg(response: Response): Option[String] = {
    if (response.getStatusCode() == http.Status.Ok.code) {
      val parsed = parseJson(response.getReader())
      parsed.extract[OPAResponse] match {
        case OPAResponse(Some(Document(errors))) =>
          if (errors.nonEmpty) {
            Some(errors.mkString("; "))
          } else {
            None
          }
        case _ =>
          Some("request denied because authorization policy is not defined")
      }
    } else {
      val parsed = parseJson(response.getReader())
      Some(parsed.extract[OPAError].toString)
    }
  }
}

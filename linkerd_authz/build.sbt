
def twitterUtil(mod: String) =
  "com.twitter" %% s"util-$mod" %  "6.40.0"

def finagle(mod: String) =
  "com.twitter" %% s"finagle-$mod" % "6.41.0"

def linkerd(mod: String) =
  "io.buoyant" %% s"linkerd-$mod" % "0.9.0"

val authzIndentifier =
  project.in(file(".")).
    settings(
      scalaVersion := "2.11.7",
      organization := "org.openpolicyagent",
      name := "linkerd-authz-identifier",
      resolvers ++= Seq(
        "twitter" at "https://maven.twttr.com",
        "local-m2" at ("file:" + Path.userHome.absolutePath + "/.m2/repository")
      ),
      libraryDependencies ++=
        finagle("http") % "provided" ::
        twitterUtil("core") % "provided" ::
        linkerd("core") % "provided" ::
        linkerd("protocol-http") % "provided" ::
        "org.json4s" %% "json4s-jackson" % "3.5.0" ::
        Nil,
      assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = false)
    )

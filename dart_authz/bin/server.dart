import 'dart:io';
import 'package:opa_api_authz_dart/opa_api_authz_dart.dart' as opa_api_authz_dart;

Future main(List<String> arguments) async {
  // Load policies into OPA server
  await opa_api_authz_dart.load_rego_policies();

  opa_api_authz_dart.pump_db();

  final HttpServer server = await HttpServer.bind(
    InternetAddress.anyIPv4, 8080,
  );

  print('Example Service listening on ' +
      '${server.address.address}:${server.port}');

  try {
    await for (HttpRequest req in server) {
      await opa_api_authz_dart.handleRequest(req);
    }
  } catch(_) {
    print("Terminating...");
    exitCode = 2;
  }
}

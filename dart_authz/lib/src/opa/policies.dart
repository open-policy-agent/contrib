import 'package:path/path.dart';
import 'package:http/http.dart' as http;
import 'package:glob/glob.dart';
import 'dart:io';
import '../helpers.dart';

String opaPolicyUrl(String host, File f) {
  return host + '/v1/policies/' + basenameWithoutExtension(f.path);
}

Future<void> load_rego_policy(FileSystemEntity f) async {
  final opaUri = getenv("OPA_URL", "http://localhost:8181");
  final file = (f as File);
  final policyUri = opaPolicyUrl(opaUri, file);
  final policy = file.readAsStringSync();

  print("Applying policy: " + file.path);

  var response = await http.put(policyUri, body: policy);
  if (response.statusCode != HttpStatus.ok) {
    print("Failed to upload policy: " + response.body);
  }
}

Future<void> load_rego_policies() async {
  final regoPolicies = Glob("**.rego");
  final regoList = regoPolicies.listSync().toList();

  await regoList.forEach((f) => load_rego_policy(f));
}

import 'dart:io';
import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/input.dart';
import '../helpers.dart';

Future<bool> request_authorized(HttpRequest req) async {
  final opaUri = getenv("OPA_URL", "http://localhost:8181");
  var input = Input();

  input.method = req.method;
  input.path = req.uri.path.substring(1).split('/');
  input.user = req.headers.value('Authorization') ?? '';

  var response = await http.post(opaUri, body: jsonEncode(input));
  if (response.statusCode != HttpStatus.ok) {
    print("Request failed with status code: " + response.statusCode.toString());
    return false;
  }

  return response.body == "true";
}
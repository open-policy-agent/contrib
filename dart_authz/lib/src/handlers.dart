import 'dart:convert';
import 'dart:io';
import 'models/cars.dart';
import 'models/mockdb.dart';
import 'opa/authorization.dart';

Future<void> handleListCars(HttpRequest req) async {
  HttpResponse resp = req.response;

  resp
    ..statusCode = HttpStatus.ok
    ..write(jsonEncode(cars));

  await resp.close();
}

Future<void> handleCarRequest(String id, HttpRequest req) async {
  HttpResponse resp = req.response;

  switch (req.method) {
    case 'GET':
      var car = cars.firstWhere((car) => car.id == id, orElse: () => null);
      if (car != null) {
        resp
          ..statusCode = HttpStatus.ok
          ..write(jsonEncode(car));
      } else {
        resp
          ..statusCode = HttpStatus.notFound
          ..write('No matching car found');
      }
      break;
    case 'DELETE':
      var numCarsBefore = cars.length;
      cars.removeWhere((car) => car.id == id);
      var numCarsAfter = cars.length;
      // Make sure the deletion has altered the list
      if (numCarsAfter == numCarsBefore) {
        resp.statusCode = HttpStatus.badRequest;
      } else {
        resp.statusCode = HttpStatus.ok;
      }
      break;
    case 'PUT':
      String content = await utf8.decoder.bind(req).join();
      var data = jsonDecode(content) as Map;
      var car = Car.fromJson(data);
      if (car == null) {
        resp
          ..statusCode = HttpStatus.unprocessableEntity
          ..write('Failed to decode Car data');
      } else {
        cars.add(car);
        resp
          ..statusCode = HttpStatus.created
          ..write('Successfully created');
      }
      break;
    default:
      resp
        ..statusCode = HttpStatus.methodNotAllowed
        ..write('Unsupported request: ${req.method}.');
      break;
  }

  await resp.close();
}

Future<void> handleCarStatusRequest(String id, HttpRequest req) async {
  HttpResponse resp = req.response;

  if (req.method == 'GET') {
    var s = status.firstWhere((car) => car.id == id, orElse: () => null);
    if (s != null) {
      resp
        ..statusCode = HttpStatus.ok
        ..write(jsonEncode(s));
    } else {
      resp
        ..statusCode = HttpStatus.notFound
        ..write('No matching car found');
    }
  } else if (req.method == 'PUT') {
    String content = await utf8.decoder.bind(req).join();
    var data = jsonDecode(content) as Map;
    var s = CarStatus.fromJson(data);
    if (s == null) {
      resp
        ..statusCode = HttpStatus.unprocessableEntity
        ..write('Failed to decode Car Status data');
    } else {
      status.add(s);
      resp
        ..statusCode = HttpStatus.created
        ..write('Successfully created');
    }
  } else {
    resp
      ..statusCode = HttpStatus.methodNotAllowed
      ..write('Unsupported request: ${req.method}.');
  }

  await req.response.close();
}

Future<void> handleRequest(HttpRequest req) async {
  HttpResponse resp = req.response;
  bool authorized = await request_authorized(req);

  if (authorized == false) {
    resp
      ..statusCode = HttpStatus.unauthorized
      ..write("Unauthorized request");
    await resp.close();
    return;
  }

  List<String> segments = req.uri.pathSegments;

  if (segments[0] != 'cars') {
    resp.statusCode = HttpStatus.badRequest;
    await resp.close();
    return;
  }

  if (segments.length == 1) {
    if (req.method == 'GET') {
      await handleListCars(req);
      return;
    }
  } else {
    var id = segments[1];

    if (segments.length == 2) {
      await handleCarRequest(id, req);
      return;
    } else if (segments.length == 3 && segments[2] == 'status') {
      await handleCarStatusRequest(id, req);
      return;
    }
  }

  resp
    ..statusCode = HttpStatus.methodNotAllowed
    ..write('Unsupported request: ${req.method}.');

  await resp.close();
}
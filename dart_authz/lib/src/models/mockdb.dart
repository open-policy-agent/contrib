import './cars.dart';

List<Car> cars = List();
List<CarStatus> status = List();

void pump_db() {
  List<Map<String, dynamic>> mockCars = [
    {
      "id": "663dc85d-2455-466c-b2e5-76691b0ce14e",
      "model": "Honda",
      "vehicle_id": "127482",
      "owner_id": "742"
    },
    {
      "id": "6c018cfa-e9c2-4169-a61b-dd3bf3bc19a7",
      "model": "Toyota",
      "vehicle_id": "19019",
      "owner_id": "742"
    },
    {
      "id": "879a273c-a8dc-41a6-9c30-2bb92288e93b",
      "model": "Ford",
      "vehicle_id": "3784312",
      "owner_id": "6928"
    },
    {
      "id": "fca3ab25-a151-4c76-b238-9aa6ee92c374",
      "model": "Honda",
      "vehicle_id": "22781",
      "owner_id": "30390"
    }
  ];

  List<Map<String, dynamic>> mockStatus = [
    {
      "id": "663dc85d-2455-466c-b2e5-76691b0ce14e",
      "position": {
        "latitude": -39.91045,
        "longitude": -161.70716
      },
      "mileage": 742,
      "speed": 90,
      "fuel": 6.42
    },
    {
      "id": "6c018cfa-e9c2-4169-a61b-dd3bf3bc19a7",
      "position": {
        "latitude": 12.77061,
        "longitude": 9.05115
      },
      "mileage": 17384,
      "speed": 62,
      "fuel": 8.9
    },
    {
      "id": "879a273c-a8dc-41a6-9c30-2bb92288e93b",
      "position": {
        "latitude": -8.86414,
        "longitude": -142.5982
      },
      "mileage": 9347,
      "speed": 45,
      "fuel": 3.1
    },
    {
      "id": "fca3ab25-a151-4c76-b238-9aa6ee92c374",
      "position": {
        "latitude": 68.86632,
        "longitude": -92.85048
      },
      "mileage": 97698,
      "speed": 50,
      "fuel": 3.22
    }
  ];

  mockCars.forEach((m) => cars.add(Car.fromJson(m)));
  mockStatus.forEach((s) => status.add(CarStatus.fromJson(s)));
}
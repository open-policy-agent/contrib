class Car {
  String id;
  String model;
  String vehicleId;
  String ownerId;

  Car({this.id, this.model, this.vehicleId, this.ownerId});

  Car.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    model = json['model'];
    vehicleId = json['vehicle_id'];
    ownerId = json['owner_id'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = Map<String, dynamic>();
    data['id'] = this.id;
    data['model'] = this.model;
    data['vehicle_id'] = this.vehicleId;
    data['owner_id'] = this.ownerId;
    return data;
  }

  @override
  String toString() {
    return 'Car[id=$id, model=$model, vehicle_id=$vehicleId, owner_id=$ownerId]';
  }
}

class CarStatus {
  String id;
  Position position;
  int mileage;
  int speed;
  double fuel;

  CarStatus({this.id, this.position, this.mileage, this.speed, this.fuel});

  CarStatus.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    position = json['position'] != null
        ? Position.fromJson(json['position'])
        : null;
    mileage = json['mileage'];
    speed = json['speed'];
    fuel = json['fuel'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = Map<String, dynamic>();
    data['id'] = this.id;
    if (this.position != null) {
      data['position'] = this.position.toJson();
    }
    data['mileage'] = this.mileage;
    data['speed'] = this.speed;
    data['fuel'] = this.fuel;
    return data;
  }

  @override
  String toString() {
    return 'CarStatus[id=$id, position=$position, mileage=$mileage, speed=$speed, fuel=$fuel]';
  }
}

class Position {
  double latitude;
  double longitude;

  Position({this.latitude, this.longitude});

  Position.fromJson(Map<String, dynamic> json) {
    latitude = json['latitude'];
    longitude = json['longitude'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = Map<String, dynamic>();
    data['latitude'] = this.latitude;
    data['longitude'] = this.longitude;
    return data;
  }

  @override
  String toString() {
    return 'Position[latitude=$latitude, longitude=$longitude]';
  }
}
class Input {
  String method;
  List<String> path;
  String user;

  Input({this.method, this.path, this.user});

  Input.fromJson(Map<String, dynamic> json) {
    method = json['method'];
    path = json['path'];
    user = json['user'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = Map<String, dynamic>();
    data['method'] = this.method;
    data['path'] = this.path;
    data['user'] = this.user;
    return data;
  }

  @override
  String toString() {
    return 'Input[method=$method, path=$path, user=$user]';
  }
}
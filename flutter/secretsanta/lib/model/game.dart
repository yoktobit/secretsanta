import 'package:secretsanta/model/status.dart';

class Game {
  final String code;
  final String name;
  final String description;
  final Status status;
  final String gifted;

  Game(this.code, this.name, this.description, this.status, this.gifted);

  factory Game.fromMap(Map<String, dynamic> json) {
    return Game(
      json['code'],
      json['name'],
      json['description'],
      json['status'],
      json['gifted'],
    );
  }

  Map<String, dynamic> toMap() {
    return {
      "code": this.code,
      "name": this.name,
      "description": this.description,
      "status": this.status,
      "gifted": this.gifted,
    };
  }
}

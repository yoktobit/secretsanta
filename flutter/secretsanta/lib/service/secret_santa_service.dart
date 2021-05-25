import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:secretsanta/model/game.dart';
import 'package:secretsanta/model/login.dart';
import 'package:secretsanta/model/login_response.dart';

class SecretSantaService {
  String baseUri = "localhost:8080";

  Future<Game> getGame() async {
    final response = await http.get(
      Uri.parse(baseUri + "/api/getGame"),
    );
    if (response.statusCode == 200) {
      return Game.fromMap(jsonDecode(response.body));
    } else {
      return null;
    }
  }

  Future<LoginResponse> login(Login login) async {
    try {
      final response = await http.post(
        Uri.parse(baseUri + "/api/loginPlayer"),
        body: jsonEncode(login.toMap()),
      );
      if (response.statusCode == 200) {
        return LoginResponse.fromMap(jsonDecode(response.body));
      } else {
        return null;
      }
    } on http.ClientException catch (exception) {
      return LoginResponse(false, exception.message);
    }
  }
}

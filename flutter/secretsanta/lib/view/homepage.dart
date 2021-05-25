import 'dart:html';

import 'package:flutter/material.dart';
import 'package:secretsanta/model/login.dart';
import 'package:secretsanta/service/secret_santa_service.dart';

class HomePage extends StatefulWidget {
  HomePage({Key key, this.title}) : super(key: key);

  final String title;
  String message = "";

  @override
  HomePageState createState() => HomePageState();
}

class HomePageState extends State<HomePage> {
  final gameCodeController = TextEditingController();
  final userController = TextEditingController();
  final passwordController = TextEditingController();
  final errorMessageController = TextEditingController();
  final secretSantaService = SecretSantaService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Text(
              errorMessageController.text,
              style: TextStyle(color: Colors.redAccent),
            ),
            ListTile(
              leading: Icon(Icons.person),
              title: TextFormField(
                controller: userController,
                decoration: InputDecoration(labelText: "Name"),
              ),
            ),
            ListTile(
              leading: Icon(Icons.vpn_key),
              title: TextFormField(
                controller: passwordController,
                obscureText: true,
                decoration: InputDecoration(labelText: "Passwort"),
              ),
            ),
            ListTile(
              leading: Icon(Icons.group),
              title: TextFormField(
                controller: gameCodeController,
                decoration: InputDecoration(labelText: "Game Code"),
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          login(context);
        },
        tooltip: 'Login',
        child: Icon(Icons.login),
      ),
    );
  }

  login(BuildContext context) async {
    final result = await secretSantaService.login(Login(
        gameCodeController.text, userController.text, passwordController.text));
    if (result.ok) {
      Navigator.pushNamed(context, "/gameAdmin");
    } else {
      errorMessageController.text = result.message;
    }
    setState(() {});
  }
}

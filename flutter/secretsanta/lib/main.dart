import 'package:flutter/material.dart';

import 'view/homepage.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: "Secret Santa's Helper",
      theme: ThemeData(
        primarySwatch: Colors.orange,
      ),
      initialRoute: "/",
      routes: {"/": (context) => HomePage(title: "Secret Santa's Helper")},
    );
  }
}

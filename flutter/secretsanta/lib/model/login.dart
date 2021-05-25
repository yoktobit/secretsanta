class Login {
  final String gameCode;
  final String user;
  final String password;

  Login(this.gameCode, this.user, this.password);

  factory Login.fromMap(Map<String, dynamic> json) {
    return Login(
      json["gameCode"],
      json["user"],
      json["password"],
    );
  }

  Map<String, dynamic> toMap() {
    return {
      "gameCode": this.gameCode,
      "user": this.user,
      "password": this.password,
    };
  }
}

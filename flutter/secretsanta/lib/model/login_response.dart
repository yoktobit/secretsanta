class LoginResponse {
  final bool ok;
  final String message;

  LoginResponse(this.ok, this.message);

  factory LoginResponse.fromMap(Map<String, dynamic> json) {
    return LoginResponse(
      json["ok"],
      json["message"],
    );
  }

  Map<String, dynamic> toMap() {
    return {
      "ok": this.ok,
      "message": this.message,
    };
  }
}

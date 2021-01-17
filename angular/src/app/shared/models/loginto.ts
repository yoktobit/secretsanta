export class LoginTo {

    private gameCode: string;
    private username: string;
    private password: string;

    constructor(gameCode: string, username: string, password: string) {
        this.gameCode = gameCode;
        this.username = username;
        this.password = password;
    }
}
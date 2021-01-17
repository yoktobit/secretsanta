export class CreateGameResponseTo {
    code: string
    link: string

    constructor(code: string, link: string) {
        this.code = code;
        this.link = link;
    }
}
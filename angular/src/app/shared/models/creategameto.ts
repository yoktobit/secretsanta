export class CreateGameTo {
    title: string
    description: string
    adminUser: string
    adminPassword: string

    constructor(title: string, description: string, adminUser: string, adminPassword: string) {

        this.title = title;
        this.description = description;
        this.adminUser = adminUser;
        this.adminPassword = adminPassword;
    }
}
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { GameCreatedComponent } from './game-created/game-created.component';
import { GameComponent } from './game/game.component';
import { LoggedInGuard } from './logged-in.guard';
import { LoginComponent } from './login/login.component';
import { ManagePlayersComponent } from './manage-players/manage-players.component';
import { NewGameComponent } from './new-game/new-game.component';
import { WelcomeComponent } from './welcome/welcome.component';

const routes: Routes = [
  {path: "", component: WelcomeComponent},
  {path: "newGame", component: NewGameComponent},
  {path: "login", component: LoginComponent},
  {path: "gameCreated", component: GameCreatedComponent},
  {path: "players", component: ManagePlayersComponent, canActivate: [LoggedInGuard]},
  {path: "game/:gameCode", component: GameComponent},
  {path: "game", component: GameComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

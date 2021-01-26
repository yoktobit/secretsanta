import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from './../environments/environment';
import { AddExceptionTo } from './shared/models/add-exception-to.model';
import { AddRemovePlayerTo } from './shared/models/add-remove-player-to.model';
import { CreateGameResponseTo } from './shared/models/create-game-response-to.model';
import { CreateGameTo } from './shared/models/create-game-to.model';
import { DrawGameResultTo } from './shared/models/draw-game-result-to.model';
import { Exception } from './shared/models/exception.model';
import { LoginResponseTo } from './shared/models/login-response-to.model';
import { LoginTo } from './shared/models/login-to.model';
import { Player } from './shared/models/player.model';
import { StatusResultTo } from './shared/models/status-result-to.model';

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  defaultOptions = { withCredentials: true };

  url = environment.backendUrl;

  constructor(private http: HttpClient) { }

  createGame(createGameTo: CreateGameTo): Observable<CreateGameResponseTo> {
    return this.http.post<CreateGameResponseTo>(this.url + "api/createNewGame", createGameTo, this.defaultOptions);
  }

  getBasicGame(gameCode: string) {
    return this.http.get(this.url + "api/game/" + gameCode)
  }

  getGame() {
    return this.http.get(this.url + "api/game", this.defaultOptions);
  }

  getPlayers() {
    return this.http.get<Player[]>(this.url + "api/players", this.defaultOptions);
  }

  getExceptions() {
    return this.http.get<Exception[]>(this.url + "api/exceptions", this.defaultOptions);
  }

  getStatus(): Observable<StatusResultTo> {
    return this.http.get<StatusResultTo>(this.url + "api/status", this.defaultOptions)
  }

  addPlayer(addPlayerTo: AddRemovePlayerTo): Observable<void> {
    return this.http.post<void>(this.url + "api/addPlayer", addPlayerTo, this.defaultOptions);
  }

  removePlayer(removePlayerTo: AddRemovePlayerTo): Observable<void> {
    return this.http.post<void>(this.url + "api/removePlayer", removePlayerTo, this.defaultOptions);
  }

  addException(addExceptionTo: AddExceptionTo) {
    return this.http.post(this.url + "api/addException", addExceptionTo, this.defaultOptions);
  }

  draw(): Observable<DrawGameResultTo> {
    return this.http.get<DrawGameResultTo>(this.url + "api/draw", this.defaultOptions);
  }

  reset() {
    return this.http.get(this.url + "api/reset", this.defaultOptions);
  }

  login(loginTo: LoginTo): Observable<LoginResponseTo> {
    return this.http.post<LoginResponseTo>(this.url + "api/loginPlayer", loginTo, this.defaultOptions);
  }

  logout() {    
    return this.http.get(this.url + "api/logout", this.defaultOptions)
  }
}

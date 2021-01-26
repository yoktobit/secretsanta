import { SelectionModel } from '@angular/cdk/collections';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { AddRemovePlayerTo } from '../shared/models/add-remove-player-to.model';
import { BackendService } from './../backend.service';
import { AddExceptionTo } from './../shared/models/add-exception-to.model';
import { Exception } from './../shared/models/exception.model';
import { Player } from './../shared/models/player.model';
import { StatusResultTo } from './../shared/models/status-result-to.model';
import { StatusService } from './../status.service';

@Component({
  selector: 'app-manage-players',
  templateUrl: './manage-players.component.html',
  styleUrls: ['./manage-players.component.css']
})
export class ManagePlayersComponent implements OnInit {

  loggedIn: boolean = false
  game: any;
  created: boolean = false;
  players: Player[] = [];
  exceptions: Exception[] = [];
  displayedColumnsPlayerList: string[] = ["select", "name", "status"]
  displayedColumns: string[] = ['nameA', 'direction', 'nameB'];
  selection = new SelectionModel<Player>(true, []);
  status: StatusResultTo;

  constructor(private router: Router, private http: HttpClient, private backend: BackendService, private statusService: StatusService) {
    this.game = this.router.getCurrentNavigation()?.extras.state;
    if (this.game) {
      this.created = true;
    } else {
      this.created = false;
    }
  }

  ngOnInit(): void {
    this.statusService.statusSubject.subscribe(status => {
      this.status = status;
      if (status.loggedIn) {
        this.loggedIn = true
        this.refreshGame();
      }
    })
  }

  baseUrl(): string {
    
    return window.location.origin;
  }

  gameUrl(): string {

    return this.baseUrl() + "/game/" + this.game?.code;
  }

  onSubmit(f: NgForm) {

    let addPlayerTo = new AddRemovePlayerTo(f.value.name);
    this.backend.addPlayer(addPlayerTo).subscribe(() => {
      this.refreshGame();
      f.resetForm();
    })
  }

  deletePlayer() {
    
    this.selection.selected.forEach((element) => {
      let removePlayerTo = new AddRemovePlayerTo(element.name);
      this.backend.removePlayer(removePlayerTo).subscribe(() => {
        this.refreshPlayers();
        this.refreshExceptions();
      })
    });
  }

  async addException() {
    
    for (const elementA of this.selection.selected) {
      for (const elementB of this.selection.selected) {
        if (elementA != elementB) {
          let addExceptionTo = new AddExceptionTo(elementA.name, elementB.name);
          await this.backend.addException(addExceptionTo).toPromise();
        }
      }
    }
    this.refreshExceptions();
  }

  getStatusText(status: string): string {
    
    switch(status) {
      case "Ready":
        return $localize `Bereit`;
      case "":
        return $localize `Registrierung steht aus`;
      default:
        return ""
    }
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    const numSelected = this.selection.selected.length;
    const numRows = this.players.length;
    return numSelected === numRows;
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  masterToggle() {
    this.isAllSelected() ?
        this.selection.clear() :
        this.players.forEach(row => this.selection.select(row));
  }

  refreshGame() {
    this.backend.getGame().subscribe(game => {
      this.game = game
      console.log(this.game)
      this.refreshPlayers();
      this.refreshExceptions();
    })
  }

  refreshPlayers() {
    this.backend.getPlayers().subscribe(players => {
      console.log("Players: " + JSON.stringify(players));
      this.players = players;
    });
  }

  refreshExceptions() {
    this.backend.getExceptions().subscribe(exceptions => {
      console.log("Exceptions: " + JSON.stringify(exceptions));
      this.exceptions = exceptions;
    });
  }

  share() {
    navigator.share({
      title: this.game.name,
      text: this.game.description,
      url: this.gameUrl()
    });
  }
}

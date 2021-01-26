import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { BackendService } from './../backend.service';
import { ConfirmResetDialogComponent } from './../confirm-reset-dialog/confirm-reset-dialog.component';
import { StatusResultTo } from './../shared/models/statusresultto';
import { StatusService } from './../status.service';

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.css']
})
export class GameComponent implements OnInit {

  game: any;
  linkedGame: any;
  loggedIn: boolean = false;
  status: StatusResultTo;

  state: string = "none";

  constructor(private route: ActivatedRoute, private router: Router, private backend: BackendService, private statusService: StatusService, private dialog: MatDialog) {
    
    this.statusService.refreshStatus();
    this.statusService.statusSubject.subscribe(status => {
      this.status = status;
    });
    this.receiveContent();
  }

  async receiveContent() {
    const gameCode = this.route.snapshot.paramMap.get("gameCode");
    const fromGameLink = gameCode != null;
    
    if (gameCode != null) {
      try {
        this.linkedGame = await this.backend.getBasicGame(gameCode).toPromise();
        this.state = "basic";
      } catch {
        this.linkedGame = null;
      }
    }
    try {
      this.game = await this.backend.getGame().toPromise();
      this.loggedIn = true;
      this.state = "detail";
    } catch {
      this.loggedIn = false;
      this.game = null;
    }
    if (this.loggedIn && fromGameLink && this.game.code != this.linkedGame.code) {
      this.state = "detailbasic";
    }
    if (this.state == "none") {
      this.router.navigate([""]);
    }
  }

  ngOnInit(): void {
  }

  async logout() {
    await this.backend.logout().toPromise();
    this.statusService.refreshStatus();
    this.router.navigate([""]);
  }

  login() {
    this.router.navigate(["/login"], {state: { "gameCode": this.linkedGame.code}})
  }
  
  startGame() {
    this.backend.draw().subscribe(() => {
      this.receiveContent();
    });
  }
  
  resetGame() {
    let dialogRef = this.dialog.open(ConfirmResetDialogComponent);
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.backend.reset().subscribe(() => {
          this.receiveContent();
          this.statusService.refreshStatus();
        })
      }
    });
  }

  getStatusText(): string {
    switch (this.game?.status) {
      case "Created":
        return $localize `Erstellt`
      case "Waiting":
        return $localize `Warte auf alle Mitwichtel`
      case "Ready":
        return $localize `Warte auf Auslosung`
      case "Drawn":
        return $localize `Ausgelost`
      default:
        return $localize `Unbekannt`
    }
  }
}

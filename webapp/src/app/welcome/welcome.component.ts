import { Component, OnInit } from '@angular/core';
import { StatusResultTo } from './../shared/models/status-result-to.model';
import { StatusService } from './../status.service';

@Component({
  selector: 'app-welcome',
  templateUrl: './welcome.component.html',
  styleUrls: ['./welcome.component.css']
})
export class WelcomeComponent implements OnInit {

  status: StatusResultTo;

  constructor(private statusService: StatusService) {
    statusService.statusSubject.subscribe(status => {
      this.status = status;
    });
  }

  ngOnInit(): void {
    this.statusService.refreshStatus();
  }

  loginButtonText(): string {
    return this.status.loggedIn ? $localize `An anderem Spiel anmelden` : $localize `An Spiel anmelden`;
  }
}

import { Component, OnInit } from '@angular/core';
import { StatusResultTo } from './../shared/models/statusresultto';
import { StatusService } from './../status.service';

@Component({
  selector: 'app-menu',
  templateUrl: './menu.component.html',
  styleUrls: ['./menu.component.css']
})
export class MenuComponent implements OnInit {

  status: StatusResultTo;

  constructor(private statusService: StatusService) { }

  ngOnInit(): void {
    this.statusService.statusSubject.subscribe(status => {
      this.status = status;
    })
  }
}

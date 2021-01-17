import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { StatusService } from './status.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  
  constructor(private statusService: StatusService, private titleService: Title) {

    this.statusService.refreshStatus();
    const title = $localize `Wichtel-Helfer`
    this.titleService.setTitle(title);
  }

  ngOnInit() {
    
  }
}

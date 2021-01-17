import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { BackendService } from './backend.service';
import { StatusResultTo } from './shared/models/statusresultto';

@Injectable({
  providedIn: 'root'
})
export class StatusService {

  private _statusSubject: BehaviorSubject<StatusResultTo> = new BehaviorSubject(new StatusResultTo(false, "", "", ""));

  constructor(private backend: BackendService) {

  }

  public get statusSubject() : BehaviorSubject<StatusResultTo> {
    return this._statusSubject;
  }  

  refreshStatus() {
    this.backend.getStatus().subscribe(status => {
      this.statusSubject.next(status);
    });
  }

}

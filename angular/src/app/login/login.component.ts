import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { BackendService } from './../backend.service';
import { LoginTo } from './../shared/models/loginto';
import { StatusService } from './../status.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  gameCode: string;

  loginForm: FormGroup;

  message: string;

  constructor(private router: Router, private formBuilder: FormBuilder, private backend: BackendService, private StatusService: StatusService) {
    this.gameCode = this.router.getCurrentNavigation()?.extras.state?.gameCode;
  }

  ngOnInit(): void {
    this.loginForm = this.formBuilder.group({
      code: [this.gameCode, Validators.required],
      username: ['', Validators.required],
      password: ['', Validators.required]
    })
    
    
  }

  onSubmit(): void {
    const loginTo = new LoginTo(this.loginForm.value.code, this.loginForm.value.username, this.loginForm.value.password);
    this.backend.login(loginTo).subscribe(result => {
      console.log(result);
      if (result.ok) {
        this.router.navigate(["/game"]);
        this.message = "";
      } else {
        this.message = result.message;
      }
      this.StatusService.refreshStatus();
    })
  }
}

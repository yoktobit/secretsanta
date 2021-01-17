import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { BackendService } from './../backend.service';
import { CreateGameTo } from './../shared/models/creategameto';

@Component({
  selector: 'app-new-game',
  templateUrl: './new-game.component.html',
  styleUrls: ['./new-game.component.css']
})
export class NewGameComponent implements OnInit {

  newGameForm: FormGroup = new FormGroup({});

  constructor(private backend: BackendService, private router: Router, private formBuilder: FormBuilder) { }

  ngOnInit(): void {
    this.newGameForm = this.formBuilder.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
      adminUser: ['', Validators.required],
      adminPassword: ['', Validators.required]
    });
  }

  onSubmit(): void {

    let createGameTo = new CreateGameTo(this.newGameForm.value.title, this.newGameForm.value.description, this.newGameForm.value.adminUser, this.newGameForm.value.adminPassword)
    this.backend.createGame(createGameTo).subscribe(createGameResponseTo  => {
      this.router.navigate(['/players'], {state: createGameResponseTo});
    }, error => {
      console.log(error);
    })
  }

}

import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { BackendService } from './../backend.service';
import { CreateGameTo } from './../shared/models/create-game-to.model';

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

  async onSubmit(): Promise<void> {

    let createGameTo = new CreateGameTo(this.newGameForm.value.title, this.newGameForm.value.description, this.newGameForm.value.adminUser, this.newGameForm.value.adminPassword)
    try {
    let createGameResponseTo = await this.backend.createGame(createGameTo);
    this.router.navigate(['/players'], {state: createGameResponseTo});
    } catch (error) {
      console.error(error);
    }
  }

}

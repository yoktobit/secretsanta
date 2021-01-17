import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GameCreatedComponent } from './game-created.component';

describe('GameCreatedComponent', () => {
  let component: GameCreatedComponent;
  let fixture: ComponentFixture<GameCreatedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ GameCreatedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(GameCreatedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfirmResetDialogComponent } from './confirm-reset-dialog.component';

describe('ConfirmResetDialogComponent', () => {
  let component: ConfirmResetDialogComponent;
  let fixture: ComponentFixture<ConfirmResetDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ConfirmResetDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfirmResetDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RawDataModalComponent } from './raw-data-modal.component';

describe('RawDataModalComponent', () => {
  let component: RawDataModalComponent;
  let fixture: ComponentFixture<RawDataModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RawDataModalComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RawDataModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

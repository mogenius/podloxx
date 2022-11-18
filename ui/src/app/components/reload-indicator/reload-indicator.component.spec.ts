import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReloadIndicatorComponent } from './reload-indicator.component';

describe('ReloadIndicatorComponent', () => {
  let component: ReloadIndicatorComponent;
  let fixture: ComponentFixture<ReloadIndicatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReloadIndicatorComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReloadIndicatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

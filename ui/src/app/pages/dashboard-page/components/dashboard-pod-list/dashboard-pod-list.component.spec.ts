import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardPodListComponent } from './dashboard-pod-list.component';

describe('DashboardPodListComponent', () => {
  let component: DashboardPodListComponent;
  let fixture: ComponentFixture<DashboardPodListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DashboardPodListComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DashboardPodListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

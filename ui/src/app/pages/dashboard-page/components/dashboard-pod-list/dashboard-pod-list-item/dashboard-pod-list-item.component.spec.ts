import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardPodListItemComponent } from './dashboard-pod-list-item.component';

describe('DashboardPodListItemComponent', () => {
  let component: DashboardPodListItemComponent;
  let fixture: ComponentFixture<DashboardPodListItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DashboardPodListItemComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DashboardPodListItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

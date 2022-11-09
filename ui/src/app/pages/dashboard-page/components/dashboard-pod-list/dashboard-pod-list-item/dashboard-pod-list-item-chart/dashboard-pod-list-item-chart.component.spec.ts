import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardPodListItemChartComponent } from './dashboard-pod-list-item-chart.component';

describe('DashboardPodListItemChartComponent', () => {
  let component: DashboardPodListItemChartComponent;
  let fixture: ComponentFixture<DashboardPodListItemChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DashboardPodListItemChartComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DashboardPodListItemChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

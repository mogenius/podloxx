import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TileOverviewComponent } from './tile-overview.component';

describe('TileOverviewComponent', () => {
  let component: TileOverviewComponent;
  let fixture: ComponentFixture<TileOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TileOverviewComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TileOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

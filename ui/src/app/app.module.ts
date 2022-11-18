import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardPageComponent } from './pages/dashboard-page/dashboard-page.component';
import { NgChartsModule } from 'ng2-charts';
import { HttpClientModule } from '@angular/common/http';
import { TopNavComponent } from './components/top-nav/top-nav.component';
import { FormsModule } from '@angular/forms';
import { formatBytes } from './pipes/formatBytes';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { RawDataModalComponent } from './components/raw-data-modal/raw-data-modal.component';
import { HashLocationStrategy, LocationStrategy } from '@angular/common';
import { HeadlineComponent } from './pages/dashboard-page/components/headline/headline.component';
import { TileOverviewComponent } from './pages/dashboard-page/components/tile-overview/tile-overview.component';
import { StatusTileComponent } from './pages/dashboard-page/components/tile-overview/status-tile/status-tile.component';
import { DashboardPodListComponent } from './pages/dashboard-page/components/dashboard-pod-list/dashboard-pod-list.component';
import { DashboardPodListItemComponent } from './pages/dashboard-page/components/dashboard-pod-list/dashboard-pod-list-item/dashboard-pod-list-item.component';
import { DashboardPodListItemChartComponent } from './pages/dashboard-page/components/dashboard-pod-list/dashboard-pod-list-item/dashboard-pod-list-item-chart/dashboard-pod-list-item-chart.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardPageComponent,
    TopNavComponent,
    formatBytes,
    RawDataModalComponent,
    HeadlineComponent,
    TileOverviewComponent,
    StatusTileComponent,
    DashboardPodListComponent,
    DashboardPodListItemComponent,
    DashboardPodListItemChartComponent
  ],
  imports: [FormsModule, HttpClientModule, BrowserModule, AppRoutingModule, NgChartsModule, NgbModule],
  providers: [{ provide: LocationStrategy, useClass: HashLocationStrategy }],
  bootstrap: [AppComponent]
})
export class AppModule {}

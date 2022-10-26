import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardPageComponent } from './pages/dashboard-page/dashboard-page.component';
import { NgChartsModule } from 'ng2-charts';
import { HttpClientModule } from '@angular/common/http';
import { ReceivedBytesGraphComponent } from './components/charts/received-bytes-graph/received-bytes-graph.component';
import { TransmitBytesGraphComponent } from './components/charts/transmit-bytes-graph/transmit-bytes-graph.component';
import { TopNavComponent } from './components/top-nav/top-nav.component';
import { FormsModule } from '@angular/forms';
import { TableListComponent } from './components/table-list/table-list.component';
import { TableListItemComponent } from './components/table-list/table-list-item/table-list-item.component';
import { formatBytes } from './pipes/formatBytes';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { RawDataModalComponent } from './components/raw-data-modal/raw-data-modal.component';
import { HashLocationStrategy, LocationStrategy } from '@angular/common';

@NgModule({
  declarations: [
    AppComponent,
    DashboardPageComponent,
    TransmitBytesGraphComponent,
    ReceivedBytesGraphComponent,
    TopNavComponent,
    TableListComponent,
    TableListItemComponent,
    formatBytes,
    RawDataModalComponent
  ],
  imports: [FormsModule, HttpClientModule, BrowserModule, AppRoutingModule, NgChartsModule, NgbModule],
  providers: [{ provide: LocationStrategy, useClass: HashLocationStrategy }],
  bootstrap: [AppComponent]
})
export class AppModule {}

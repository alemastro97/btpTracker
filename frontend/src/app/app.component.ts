import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
/* Core Grid CSS */
import 'ag-grid-community/styles/ag-grid.css';
/* Quartz Theme Specific CSS */
import 'ag-grid-community/styles/ag-theme-quartz.css';
import { RouterOutlet } from '@angular/router';
import { AgGridAngular } from 'ag-grid-angular';
import { ColDef } from 'ag-grid-community';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { ScaleType } from '@swimlane/ngx-charts';

// @ts-nocheck
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, AgGridAngular, NgxChartsModule],
  template: `
  <div class="overflow-x-auto text-white absolute top-0 left-0 right-0 bottom-0 bg-[#1f2836]">
      <!-- The AG Grid component -->
<div class="w-full h-full flex">
    <div class="w-[50%] h-full">
      <ag-grid-angular
            style="width: 100%; height: 100%;"
            [class]="themeClass"
            [rowData]="data"
            [columnDefs]="colDefs"
            (cellClicked)="onCellClicked($event)"
          >
          </ag-grid-angular>
    </div>
        <div class="w-[50%] h-[90%]">

          <!-- [view]="null"
          [curve]="curve"
          [wrapTicks]="wrapTicks" -->
<ngx-charts-line-chart

        
        [animations]="animations"
        
        [autoScale]="true"
        [showGridLines]="true"
        [rangeFillOpacity]="3"
        [roundDomains]="true"
        [yAxisTickFormatting]="axisFormat"

  [legend]="legend"
  scaleType="linear"
  [timeline]="false"

  [gradient]="true"
  [showXAxisLabel]="showXAxisLabel"
  [showYAxisLabel]="showYAxisLabel"
  [xAxis]="xAxis"
  [yAxis]="yAxis"
  [schemeType]="schemeType"
  [scheme]=" {
    name: 'nightLights',
    selectable: false,
    group: scaling,
    domain: [
      '#4e31a5',
      '#9c25a7',
      '#3065ab',
      '#57468b',
      '#904497',
      '#46648b',
      '#32118d',
      '#a00fb3',
      '#1052a2',
      '#6e51bd',
      '#b63cc3',
      '#6c97cb',
      '#8671c1',
      '#b455be',
      '#7496c3'
    ]
  }
"
  [xAxisLabel]="xAxisLabel"
  [yAxisLabel]="yAxisLabel"
  [timeline]="timeline"
  [results]="dataPlot"
  >
</ngx-charts-line-chart>
        </div>
        <style lang="sass">
          ::ng-deep .ngx-charts {
            text{
                fill: #747c8c!important;
            }

            ::ng-deep .ngx-charts .gridline-path {
                stroke: red !important;
            }

          }
          

        </style>
</div>
  </div>
  `,
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  [x: string]: any;
  data: any = [];
  dataPlot: any = [];
  themeClass = "ag-theme-quartz-auto-dark";
  colDefs: ColDef[] = [
    { field: "ISIN", cellDataType: 'text' },
    { field: "Description", cellDataType: 'text' },
    { field: "Last", cellDataType: 'number' },
    { field: "Cedola", cellDataType: 'number' },
    { field: "Expiration", cellDataType: 'date' }
  ];

  // multi: any[];
  // view: any[] = [600, 300];

  // options
  legend: boolean = true;
  showLabels: boolean = true;
  animations: boolean = true;
  xAxis: boolean = true;
  yAxis: boolean = true;
  showYAxisLabel: boolean = true;
  showXAxisLabel: boolean = true;
  xAxisLabel: string = 'Year';
  yAxisLabel: string = 'Population';
  timeline: boolean = true;


  scaling = ScaleType.Linear
  schemeType = ScaleType.Linear

  axisFormat = (value: number) => {
    console.log(value % 1 === 0);
    if (value % 1 === 0) {
      return value.toLocaleString();
    } else {
      return value.toString();
    }
  }

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.http.get('http://localhost:8080/getRTData').subscribe(response => {

      this.data = (response as any[]).map(({
        ISIN, Description, Last, Cedola, Expiration
      }) => ({
        ISIN,
        Description,
        Last: Number.isNaN(parseFloat(Last.replace(',', '.'))) ? Infinity : parseFloat(Last.replace(',', '.')),
        Cedola: Number.isNaN(parseFloat(Cedola.replace(',', '.'))) ? Infinity : parseFloat(Cedola.replace(',', '.')),
        Expiration: new Date(Expiration.replace(',', '.'))
      })).filter(({ ISIN }) => ISIN);

    })

  }

  onCellClicked(event: any) {
    // handle cell click event here
    // alert(event.data.ISIN);
      this.http.get(`http://localhost:8080/getBTPData?id=${event.data.ISIN}`).subscribe(response => {

      // this.data = (response as any[])
      this.dataPlot = [{name:event.data.ISIN, series:(response as any[]).map(({name,value}) => ({value: parseFloat(value.replace(",", '.')),name}))}];
      console.log(this.dataPlot)
      // .map(({
      //   ISIN, Description, Last, Cedola, Expiration
      // }) => ({
      //   ISIN,
      //   Description,
      //   Last: Number.isNaN(parseFloat(Last.replace(',', '.'))) ? Infinity : parseFloat(Last.replace(',', '.')),
      //   Cedola: Number.isNaN(parseFloat(Cedola.replace(',', '.'))) ? Infinity : parseFloat(Cedola.replace(',', '.')),
      //   Expiration: new Date(Expiration.replace(',', '.'))
      // })).filter(({ ISIN }) => ISIN);

    })
  }
}
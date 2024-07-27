import {Injectable, signal} from '@angular/core';
import {LocalStorageService} from "@services/local-storage.service";

export enum ModeType {
  DARK = 'dark',
  LIGHT = 'light',
}
@Injectable({
  providedIn: 'root'
})
export class ModeService {
  private key = "current_mode";

  modeSignal = signal<ModeType>(ModeType.LIGHT)


  constructor(private localStorageService: LocalStorageService) {
    this.modeSignal.update(value => value = this.getMode())
  }

  updateMode(mode: ModeType) {
    this.modeSignal.update(value => value = mode)
    this.setMode(mode);
  }

  public getMode(): ModeType {
    let mode = this.localStorageService.get('current_mode');
    return mode === ModeType.LIGHT ? ModeType.LIGHT : ModeType.DARK;
  }

  public setMode(mode: ModeType) {
    this.localStorageService.set('current_mode', mode)
  }
}

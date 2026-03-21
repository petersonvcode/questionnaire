import { MainPage } from "./pages/mainPage";

export class App {
  private readonly mainPage: MainPage;


  constructor() {
    this.mainPage = new MainPage();

    this.mainPage.show();
  }
}

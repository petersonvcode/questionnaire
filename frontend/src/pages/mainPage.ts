import { Page, showElement, hideElement } from "./page";
import { AnsweringPage } from "./answeringPage";
import { getQuestionsWithTags } from "../questions";
import { RandomQuestionPage } from "./randomQuestionPage";
import { TestPage } from "./testsPage";
import { ThemesQuestionPage } from "./themesPage";
import { ResultsPage } from "./resultsPage";

export class MainPage implements Page {
  // Menu options
  private readonly testPage: TestPage;
  private readonly themesQuestionPage: ThemesQuestionPage;
  private readonly randomQuestionPage: RandomQuestionPage;

  // Intermidiate and final pages
  private readonly answeringPage: AnsweringPage;
  private readonly resultsPage: ResultsPage;

  private readonly pageElement: HTMLElement;

  constructor() {
    this.resultsPage = new ResultsPage();
    this.randomQuestionPage = new RandomQuestionPage();
    this.answeringPage = new AnsweringPage((questions, answersIndexes) => {
      this.hideAllPages();
      this.resultsPage.show();
      this.resultsPage.loadQuestions(questions, answersIndexes);
    });
    this.themesQuestionPage = new ThemesQuestionPage(async (tags, count) => {
      this.hideAllPages();
      this.answeringPage.show();
      const tagIds = tags.map(t => t.id);
      const questions = await getQuestionsWithTags(tagIds, count);
      this.answeringPage.loadQuestions(questions);
    });
    this.testPage = new TestPage();
    
    const pageElement = document.getElementById('main-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;

    // Main page buttons
    document.getElementById('test-btn')
      ?.addEventListener('click', () => {
        this.hideAllPages();
        this.testPage.show();
      });
    document.getElementById('themes-btn')
      ?.addEventListener('click', () => {
        this.hideAllPages();
        this.themesQuestionPage.show();
      });
    document.getElementById('random-question-btn')
      ?.addEventListener('click', () => {
        this.hideAllPages();
        this.randomQuestionPage.show();
      });

    // Go to main page buttons
    const goToMainBtns = document.getElementsByClassName('go-to-main-btn');
    for (const btn of goToMainBtns)
      btn.addEventListener('click', () => {
        this.hideAllPages();
        this.show();
      });
  }

  show = () => showElement(this.pageElement)

  hide = () => hideElement(this.pageElement)

  private hideAllPages(): void {
    this.hide();
    this.testPage.hide();
    this.themesQuestionPage.hide();
    this.randomQuestionPage.hide();
    this.answeringPage.hide();
    this.resultsPage.hide();
  }
}
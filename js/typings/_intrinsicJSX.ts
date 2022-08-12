/**
 * This file provides typings for built-in JSX elements. Typescript assumes no
 * elements are defined by default.
 */

/* eslint no-unused-vars: ["off"] */

declare global {
  function createElement(
      tag: string,
      attrs: {[key: string]: string},
      ...children: (Node | string)[]
  ): HTMLElement;

  namespace JSX {
    interface IntrinsicElements {
      div: any;
      span: any;
      img: any;
      a: any;
      h1: any;
      h2: any;
      h3: any;
      h4: any;
      h5: any;
      h6: any;
      p: any;
      option: any;
      sup: any;
      table: any;
      tbody: any;
      tr: any;
      th: any;
      td: any;
      col: any;
      code: any;
    }
  }
}

// Export {} to trick Typescript into thinking this is a module
// https://stackoverflow.com/a/59499895/7481517
export {};

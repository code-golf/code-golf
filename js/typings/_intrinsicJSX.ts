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
      a:        any;
      aside:    any;
      button:   any;
      code:     any;
      col:      any;
      div:      any;
      h1:       any;
      h2:       any;
      h3:       any;
      h4:       any;
      h5:       any;
      h6:       any;
      header:   any;
      img:      any;
      input:    any;
      label:    any;
      li:       any;
      optgroup: any;
      option:   any;
      p:        any;
      section:  any;
      select:   any;
      span:     any;
      sup:      any;
      svg:      any;
      table:    any;
      tbody:    any;
      td:       any;
      th:       any;
      thead:    any;
      time:     any;
      tr:       any;
      use:      any;
    }
  }
}

// Export {} to trick Typescript into thinking this is a module
// https://stackoverflow.com/a/59499895/7481517
export {};

function guess(char) {
  const currentGuess = document.getElementById("current");
  const gameover = document.getElementById("gameover");
  if (char && currentGuess && !gameover) {
    if (char === "ENTER") {
      const word = Array.from(currentGuess.children).reduce(
        (word, letter) => word + letter.innerHTML,
        ""
      );
      if (word.length === 5) {
        fetch("/guess", {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
          },
          body: new URLSearchParams({ word: word }),
        }).then((response) => {
          if (response.status === 200) {
            const gordle = document.getElementById("gordle");
            response.text().then((text) => (gordle.outerHTML = text));

            setTimeout(() => {
              const error = document.getElementById("error");
              if (error) {
                error.classList.add("hide");
              }
            }, 1000);
          }
        });
      }
    } else if (char === "BACKSPACE") {
      const lastChar = Array.from(currentGuess.children).findLast(
        (letter) => letter.innerHTML !== ""
      );
      if (lastChar) {
        lastChar.innerHTML = "";
      }
    } else if (
      char.length === 1 &&
      char.charCodeAt() >= 65 &&
      char.charCodeAt() <= 90
    ) {
      const currentLetter = Array.from(currentGuess.children).find(
        (letter) => letter.innerHTML === ""
      );
      if (currentLetter) {
        currentLetter.innerHTML = char;
      }
    }
  }
}

document.addEventListener("keydown", (event) => guess(event.key.toUpperCase()));

function reset() {
  fetch("/reset").then((response) => {
    if (response.status === 200) {
      const gordle = document.getElementById("gordle");
      response.text().then((text) => (gordle.outerHTML = text));
    }
  });
}

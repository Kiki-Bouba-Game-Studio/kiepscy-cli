const videoDivs = document.querySelectorAll(".e-tab-content");
const videos = [];
videoDivs.forEach(videoDiv => {
  let videoData = {};
  videoData.url = videoDiv.getAttribute('data-video-url');
  videoData.title = videoDiv.getAttribute('data-video-title');
  videoData.duration = videoDiv.getAttribute('data-video-duration');
  videoData.type = videoDiv.getAttribute('data-video-type');
  videoData.thumbnail = document.getElementById(videoDiv.getAttribute("aria-labelledby")).querySelector("img").src
  videos.push(videoData);
});

var jsonString = JSON.stringify(videos, null, 2);
var modifiedString = `,${jsonString.slice(1, -1)}`
console.log(modifiedString)

// wikipedia seasons data scraper
const tablesContainer = document.querySelector(".mw-content-ltr");
const seasonTitlesContainers = tablesContainer.querySelectorAll(".mw-heading")
const episodesTablesContainers = Array.from(tablesContainer.querySelectorAll(".wikitable")).slice(0, -1)

const seasons = []

seasonTitlesContainers.forEach((seasonTitleContainer, index) => {
  const seasonTitle = seasonTitleContainer.querySelector("h2").textContent;
  const table = episodesTablesContainers[index]
  if (!table) return
  const tableRows = table.querySelectorAll("tbody > tr");
  const episodes = []
  tableRows.forEach(row => {
    const cells = row.querySelectorAll("td")
    if (cells.length > 0) {
      const episodeNumber = cells[0].innerText

      // table changes at 10th season
      let i = 2
      if (index >= 10)
        i = 1
      const episodeTitle = cells[i].innerText
      const episode = {
        number: episodeNumber,
        title: episodeTitle,
      }
      episodes.push(episode)
    }
  })
  const season = {
    title: seasonTitle,
    episodes: episodes,
  }
  seasons.push(season)
})

console.log(seasons)

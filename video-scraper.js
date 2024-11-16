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

console.log(videos)


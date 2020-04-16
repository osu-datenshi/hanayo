(function() {
var reviews = [
    {
    name: "Akinya",
    avatar:  "https://a.ainu.pw/1300",
    link: "https://ainu.pw/u/1300",
    review: "this server is lovely, with friendly people but like to tease the others"
    },
    {
    name: "Zendaya",
    avatar:  "https://a.ainu.pw/1188",
    link: "https://ainu.pw/u/1188",
    review: "Lmaoo they asking for reviews for the homepage put me there or u ban"
    },
    {
    name: "Cherry",
    avatar:  "https://a.ainu.pw/1028",
    link: "https://ainu.pw/u/1028",
    review: "I want say thanks to Aoba (Server Owner), he always was help me in my beginnings and ... One more thanks <3 Thanks to Ainu Community."
    },
    {
    name: "Pepega",
    avatar: "https://a.ainu.pw/1023",
    link: "https://ainu.pw/u/1023",
    review: "most fun i've had playing osu has been on this server would recommend"
    },
    {
    name: "KeepSOBP",
    avatar: "https://a.ainu.pw/1329",
    link: "https://ainu.pw/u/1329",
    review: "Please register ainu so you can get country rank #1 faster OOF"
    },
    {
    name: "Phirse",
    avatar: "https://a.ainu.pw/1356",
    link: "https://ainu.pw/u/1356",
    review: "hi"
    },
    {
    name: "Vaxxine",
    avatar:  "https://a.ainu.pw/1663",
    link: "https://ainu.pw/u/1663",
    review: "This is a really good server I hope there will be more players in the future :))"
    },
    {
    name: "Sonic",
    avatar:  "https://a.ainu.pw/1649",
    link: "https://ainu.pw/u/1649",
    review: "Thanks to Ainu, the police started a missing search!"
    }
];
var review = reviews[Math.floor(Math.random() * 8)];
document.getElementById("review").innerHTML =
    '<a href="' + review.link + '" class="reivew-avatar" style="background-image: url(' + review.avatar + ')"></a>' +
    '<div class="review-content">' +
    '<a href="' + review.link + '" class="text-h3">' + review.name + '</a>' +
    '<span class="review-text">' + review.review  + '</span>'
    '</div>';
})();

(function() {
var quotes = [
    {
    name: "LifelessApple",
    avatar:  "https://a.ainu.pw/1264",
    link: "https://ainu.pw/u/1264",
    review: "good server I like better than other server :)"
    },
    {
    name: "0249",
    avatar:  "https://a.ainu.pw/1357",
    link: "https://ainu.pw/u/1357",
    review: "ainu very good and epic and stuff yes"
    },
    {
    name: "Doppler",
    avatar:  "https://a.ainu.pw/1367",
    link: "https://ainu.pw/u/1367",
    review: "Best server you'll ever find"
    },
    {
    name: "xxdnull",
    avatar:  "https://a.ainu.pw/1178",
    link: "https://ainu.pw/u/1178",
    review: "welcome to ainu please like and subscribe then comment down below"
    },
    {
    name: "Kiyushi",
    avatar: "https://a.ainu.pw/1538",
    link: "https://ainu.pw/u/1538",
    review: "<3"
    },
    {
    name: "susanoo",
    avatar: "https://a.ainu.pw/1627",
    link: "https://ainu.pw/u/1627",
    review: "good server keep up"
    },
    {
    name: "TBG",
    avatar: "https://a.ainu.pw/1520",
    link: "https://ainu.pw/u/1520",
    review: "good server keep up"
    },
    {
    name: "oh no",
    avatar:  "https://a.ainu.pw/1669",
    link: "https://ainu.pw/u/1669",
    review: "ok now this server is epic"
    }
];
var quote = quotes[Math.floor(Math.random() * 8)];
document.getElementById("review-row-3").innerHTML =
    '<a href="' + quote.link + '" class="reivew-avatar" style="background-image: url(' + quote.avatar + ')"></a>' +
    '<div class="review-content">' +
    '<a href="' + quote.link + '" class="text-h3">' + quote.name + '</a>' +
    '<span class="review-text">' + quote.review  + '</span>'
    '</div>';
})();

(function() {
var quotes = [
    {
    name: "chiwawa",
    avatar:  "https://a.ainu.pw/1194",
    link: "https://ainu.pw/u/1194",
    review: "ainu epic join now plox"
    },
    {
    name: "Aoffy",
    avatar:  "https://a.ainu.pw/1007",
    link: "https://ainu.pw/u/1007",
    review: "The food is amazing, air conditioner's great"
    },
    {
    name: "Mahmood",
    avatar:  "https://a.ainu.pw/1101",
    link: "https://ainu.pw/u/1101",
    review: "Bad maps get unranked pog"
    },
    {
    name: "leina",
    avatar:  "https://a.ainu.pw/1447",
    link: "https://ainu.pw/u/1447",
    review: "epic server with epic people, should totally join 11/10"
    },
    {
    name: "Rulehd",
    avatar: "https://a.ainu.pw/1038",
    link: "https://ainu.pw/u/1038",
    review: "best server."
    },
    {
    name: "Hakura",
    avatar: "https://a.ainu.pw/1231",
    link: "https://ainu.pw/u/1231",
    review: "Tangmo with a gun (jk btw it's nice server : )"
    },
    {
    name: "__",
    avatar: "https://a.ainu.pw/1206",
    link: "https://ainu.pw/u/1206",
    review: "great server, helped me get my life support back"
    }
];
var quote = quotes[Math.floor(Math.random() * 7)];
document.getElementById("review-row-2").innerHTML =
    '<a href="' + quote.link + '" class="reivew-avatar" style="background-image: url(' + quote.avatar + ')"></a>' +
    '<div class="review-content">' +
    '<a href="' + quote.link + '" class="text-h3">' + quote.name + '</a>' +
    '<span class="review-text">' + quote.review  + '</span>'
    '</div>';
})();

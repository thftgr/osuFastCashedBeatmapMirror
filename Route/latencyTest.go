package Route

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func LatencyTest(c echo.Context) (err error) {
	return c.HTML(http.StatusOK, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/chart.min.js" integrity="sha512-QSkVNOCYLtj73J4hbmVoOV6KVZuMluZlioC+trLpewV8qMjsWqlIQvkn1KGX2StWvPMdWGBqim1xlC8krl1EKQ=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/chart.esm.js" integrity="sha512-jUlTTDoq6IvZiinGFQetLcklithBTp8sVUkkUBEYQvYd3hwMuCTd59kAzVpJwvRTmZ2palO++nX+vKC+cK9lqg=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/chart.esm.min.js" integrity="sha512-ESlgC6ZyzTZdmD7XoDrXJNOqoIugH+FNKg8nAk8sa3cZfFALiV+lo5xoia649oyygwnkKsdUdPAJ+puqGbOs+g=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/chart.js" integrity="sha512-Lii3WMtgA0C0qmmkdCpsG0Gjr6M0ajRyQRQSbTF6BsrVh/nhZdHpVZ76iMIPvQwz1eoXC3DmAg9K51qT5/dEVg=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/helpers.esm.js" integrity="sha512-DfyVRytIoS7LXOzmxSD4//oV81uwh1xV/EuK/xeh5gVJycOdBj+XTl7jeM6bcy7jiBIabR/9S2uyRL3oKKustw=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.7.1/helpers.esm.min.js" integrity="sha512-urWBnIv+F027G24xDNigjxvIuwnWlWy94W2yx77VkISKLzKSohOKOubMDhtEF6LZcEH7gctmNSpxDqIW/zMmUg=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>

</head>
<body>
<div>
    <!--차트가 그려질 부분-->
    <canvas id="myChart"></canvas>
</div>
<script type="text/javascript">
    const apiserver = [
        "ko.nerinyan.moe",
        "ko2.nerinyan.moe",
        "eu.nerinyan.moe",
        "us.nerinyan.moe",
        "rus.nerinyan.moe",
    ]
    const promises = [];

    const context = document.getElementById('myChart').getContext('2d');
    const DATA_COUNT = 50;
    const labels = [];
    for (let i = 0; i < DATA_COUNT; ++i) {
        labels.push(i.toString());
    }
    const data = [];
    const datasets = [];

    const myChart = new Chart(context, {
        type: 'line', // 차트의 형태
        data: {
            labels: labels,
            datasets: datasets
        },
        options: {
            responsive: true,
            plugins: {
                title: {
                    display: true,
                    text: 'Chart.js Line Chart - Cubic interpolation mode'
                },
            },
            interaction: {
                intersect: false,
            },
            scales: {
                x: {
                    display: true,
                    title: {
                        // display: true
                    }
                },
                y: {
                    display: true,
                    title: {
                        display: true,
                        text: 'latency(ms)'
                    },
                    suggestedMin: 0,
                    suggestedMax: 200
                }
            }
        }
    });
    apiserver.map((v, i) => {
        data.push([])

        datasets.push({
            label: v,
            data: data[i],
            borderColor: "#"+Math.floor(Math.random() * 16777215).toString(16),
            fill: false,
            cubicInterpolationMode: 'monotone',
            tension: 0.2
        })
        new Promise(async () => {
            while (1) {
                const st = new Date().getTime()

                await fetch("https://"+v+"/health").then(async res => {
                    const et = new Date().getTime()
                    if (res.ok) data[i].push((et - st))

                    if (data[i].length > DATA_COUNT) {
                        data[i].shift()
                        // console.log(data[i].length)
                    }

                    await new Promise(r => setTimeout(r, 1000));
                }).then(()=> myChart.update('none')).catch(async (err) => {
                    //fetch 에러 발생시 1초 대기
                    await new Promise(r => setTimeout(r, 1000));
                })
            }
        }).then().catch(err => console.log(err))

    })

</script>

</body>
</html>
`)

}

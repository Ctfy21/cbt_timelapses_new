const axios = require('axios')

export function downloadImage(room, camera, startDate, endDate){
    const url = `http://192.168.42.119:5000/download/${room}/${camera}/timelapses/output_${startDate.toISOString().slice(0,10)}_00-00-00_to_${endDate.toISOString().slice(0,10)}_00-00-00.mp4`
    console.log(url)
    axios({
        url: url,
        method: 'GET',
        responseType: 'blob', // important
    }).then((response) => {
        // create file link in browser's memory
        const href = URL.createObjectURL(response.data);
        // create "a" HTML element with href to file & click
        const link = document.createElement('a');
        link.href = href;
        link.setAttribute('download', `timelapse_${room}_${camera}_${startDate.toISOString().slice(0,10)}_${endDate.toISOString().slice(0,10)}`); //or any other extension
        document.body.appendChild(link);
        link.click();
        // clean up "a" element & remove ObjectURL
        document.body.removeChild(link);
        URL.revokeObjectURL(href);
    });
}
module.exports = {
    apps: [
        {
            name: "homelabs-service",
            script: "/root/devlabs/tutitoos/backend/homelabs-service/bin/app",
            max_memory_restart: "1G",
            env:  {
                GOMAXPROCS: "1"
            }
        }
    ]
};

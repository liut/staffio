#web: rerun -goexec=${GOEXEC} --build --ignore fe -watch ./pkg --rundir tmp github.com/liut/staffio web --fs local
demo: rerun -goexec=${GOEXEC} --build --ignore fe -watch ./pkg --rundir tmp github.com/liut/staffio/cmd/staffio-demo
fe: ./node_modules/.bin/gulp watch

import { FunctionalComponent, h } from 'preact';
import { useEffect, useRef, useState } from 'preact/hooks';
import feedPB from '../../protobuf/feed_pb';
import style from './style.css';


const Home: FunctionalComponent = () => {
    const [address, setAddress] = useState('ws://localhost:32001/ws')
    const [connectionState, setConnectionState] = useState<number | undefined>(undefined)
    const [velocity, setVelocity] = useState<number[] | undefined>(undefined)
    const [position, setPosition] = useState<number[] | undefined>(undefined)
    const [fsu, setFSU] = useState<number[] | undefined>(undefined)
    const websocketRef = useRef<WebSocket | undefined>()
    const dataReceived = useRef(0)
    const [dataRate, setDataRate] = useState(0)
    const canvasRef = useRef<HTMLCanvasElement>()
    const drawQueue = useRef<number[][]>([])

    useEffect(() => {
        const ctx = canvasRef.current.getContext('2d')
        if (!ctx) {
            return;
        }

        ctx.strokeStyle = 'red'
        ctx.lineWidth = 1
    }, [])

    useEffect(() => {
        const transformPoint = (ctx: CanvasRenderingContext2D, [x, y]: [number, number]) => {
            x *= 0.1
            y *= 0.1
            x += ctx.canvas.width / 2
            y += ctx.canvas.height / 2
            return [Math.round(x), Math.round(y)]
        }

        const func = () => {
            window.requestAnimationFrame(func)

            if (!drawQueue.current.length) {
                return;
            }

            const ctx = canvasRef.current.getContext('2d')
            if (!ctx) {
                return;
            }

            const initialPoint = transformPoint(ctx, drawQueue.current[0] as [number, number])
            ctx.moveTo(initialPoint[0], initialPoint[1])
            for (let i = 1; i < drawQueue.current.length; i++) {
                const point = transformPoint(ctx, drawQueue.current[i] as [number, number])
                ctx.lineTo(point[0], point[1])
            }
            ctx.stroke()
            drawQueue.current = [drawQueue.current[drawQueue.current.length - 1]]
        }
        window.requestAnimationFrame(func)
    }, [])

    useEffect(() => {
        let prevDataReceived = 0
        let prevTime = 0
        const interval = setInterval(() => {
            const timeNow = +new Date()
            const timeDelta = (timeNow - prevTime) / 1000
            setDataRate((dataReceived.current - prevDataReceived) / timeDelta)
            prevTime = timeNow
            prevDataReceived = dataReceived.current
        }, 1000)
        return () => {
            clearInterval(interval)
        }
    }, [])

    const onButtonClick = () => {
        if (websocketRef.current) {
            websocketRef.current.close()
            websocketRef.current = undefined;
        }

        websocketRef.current = new WebSocket(address)
        websocketRef.current.binaryType = 'arraybuffer'
        setConnectionState(websocketRef.current.readyState)
        websocketRef.current.onopen = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
        }
        websocketRef.current.onerror = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
        }
        websocketRef.current.onclose = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
            websocketRef.current = undefined
        }
        websocketRef.current.onmessage = (e) => {
            dataReceived.current += e.data.byteLength;
            const pmove = feedPB.PMove.deserializeBinary(e.data);
            const newPosition = pmove.getPositionList();
            drawQueue.current.push(newPosition)
            setVelocity(pmove.getVelocityList())
            setPosition(newPosition)
            setFSU(pmove.getFsuList())
        }
    }

    const onAddressInput: h.JSX.GenericEventHandler<HTMLInputElement> = (e) => {
        if (e.target instanceof HTMLInputElement) {
            setAddress(e.target.value)
        }
    }

    const connectionStateText = (() => {
        if (typeof connectionState !== 'number') {
            return 'Idle'
        }
        switch (connectionState) {
            case WebSocket.CONNECTING:
                return 'Connecting';
            case WebSocket.OPEN:
                return 'Open';
            case WebSocket.CLOSING:
                return 'Closing';
            case WebSocket.CLOSED:
                return 'Closed';
            default:
                return 'Unknown';
        }
    })()

    return (
        <div class={style.home}>
            <h1>Home</h1>
            <p>This is the Home component.</p>
            <label>
                Address: <input type="text" value={address} onInput={onAddressInput} />
            </label>
            <div>Connection state: {connectionStateText}</div>
            <div>Downlink data rate: {(dataRate / 1000).toFixed(3)} KB/s</div>
            <button onClick={onButtonClick}>Connect</button>

            {velocity ? <div>Velocity: {velocity[0].toFixed(6)} {velocity[1].toFixed(6)} {velocity[2].toFixed(6)}</div> : null}
            {position ? <div>Position: {position[0].toFixed(6)} {position[1].toFixed(6)} {position[2].toFixed(6)}</div> : null}
            {fsu ? <div>FSU: {fsu[0].toFixed(6)} {fsu[1].toFixed(6)} {fsu[2].toFixed(6)}</div> : null}

            <div>
                <canvas ref={canvasRef} width={500} height={500} style="border: 1px solid" />
            </div>
        </div>
    );
};

export default Home;

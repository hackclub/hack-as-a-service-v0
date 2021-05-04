import Meta from '@hackclub/meta'

export default {
  repository: 'https://github.com/hackclub/hack-as-a-service',
  titleSuffix: ' – HaaS Docs',
  logo: (
    <>
      <img
        src="https://cloud-jo17wa721-hack-club-bot.vercel.app/0haas.png"
        style={{
          height: '20px',
          marginRight: '4px',
          transform: 'translateY(-1px)',
        }}
      />
      <span className="mr-1 font-extrabold hidden md:inline">
        Hack As A Service
      </span>
      <span className="text-gray-600 font-normal hidden md:inline">
        Documentation
      </span>
    </>
  ),
  head: (
    <Meta
      name="Hack as a Service by Hack Club" // site name
      description="How to use Hack Club's compute as a service." // page description
      image="https://workshop-cards.hackclub.com/Hack%20as%20a%20Service%20Docs.png?theme=light&md=1&fontSize=275px&caption=&images=&images=" // large summary card image URL
      color="#ec3750" // theme color
      manifest="/site.webmanifest" // link to site manifest
    />
  ),
  search: true,
  prevLinks: true,
  nextLinks: true,
  footer: true,
  footerEditOnGitHubLink: false,
  footerText: <>MIT {new Date().getFullYear()} © Hack Club.</>,
}
